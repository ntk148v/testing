import asyncio
from typing import Optional
from contextlib import AsyncExitStack

from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client
from google import genai
from google.genai import types
from dotenv import load_dotenv

load_dotenv()  # load environment variables from .env
model = "gemini-2.0-flash"


class MCPClient:
    def __init__(self):
        # Initialize session and client objects
        self.session: Optional[ClientSession] = None
        self.exit_stack = AsyncExitStack()
        self.gemini = genai.Client()

    async def connect_to_server(self, server_script_path: str):
        """Connect to an MCP server

        Args:
            server_script_path: Path to the server script (.py or .js)
        """
        is_python = server_script_path.endswith(".py")
        is_js = server_script_path.endswith(".js")
        if not (is_python or is_js):
            raise ValueError("Server script must be a .py or .js file")

        command = "python" if is_python else "node"
        server_params = StdioServerParameters(
            command=command, args=[server_script_path], env=None
        )

        stdio_transport = await self.exit_stack.enter_async_context(
            stdio_client(server_params)
        )
        self.stdio, self.write = stdio_transport
        self.session = await self.exit_stack.enter_async_context(
            ClientSession(self.stdio, self.write)
        )

        await self.session.initialize()

        # List available tools
        response = await self.session.list_tools()
        tools = response.tools
        print("\nConnected to server with tools:", [tool.name for tool in tools])

    async def process_query(self, query: str) -> str:
        """Process a query using Claude and available tools"""
        contents = [types.Content(role="user", parts=[types.Part(text=query)])]

        # 1. -- Get Tools from Session and convert to Gemini Tool objects ---
        response = await self.session.list_tools()
        tools = types.Tool(
            function_declarations=[
                {
                    "name": tool.name,
                    "description": tool.description,
                    "parameters": tool.inputSchema,
                }
                for tool in response.tools
            ]
        )

        # --- 2. Initial Request with user prompt and function declarations ---
        response = await self.gemini.aio.models.generate_content(
            model=model,
            contents=contents,
            config=types.GenerateContentConfig(
                temperature=0,
                tools=[tools],
            ),  # Example other config
        )

        # --- 3. Append initial response to contents ---
        contents.append(response.candidates[0].content)

        # --- 4. Tool Calling Loop ---
        turn_count = 0
        max_tool_turns = 5

        while response.function_calls and turn_count < max_tool_turns:
            turn_count += 1
            tool_response_parts: List[types.Part] = []

            # --- 4.1 Process all function calls in order and return in this turn ---
            for fc_part in response.function_calls:
                tool_name = fc_part.name
                args = fc_part.args or {}  # Ensure args is a dict
                print(f"Attempting to call MCP tool: '{tool_name}' with args: {args}")

                tool_response: dict
                try:
                    # Call the session's tool executor
                    tool_result = await self.session.call_tool(tool_name, args)
                    print(f"MCP tool '{tool_name}' executed successfully.")
                    if tool_result.isError:
                        tool_response = {"error": tool_result.content[0].text}
                    else:
                        tool_response = {"result": tool_result.content[0].text}
                except Exception as e:
                    tool_response = {
                        "error": f"Tool execution failed: {type(e).__name__}: {e}"
                    }

                # Prepare FunctionResponse Part
                tool_response_parts.append(
                    types.Part.from_function_response(
                        name=tool_name, response=tool_response
                    )
                )

            # --- 4.2 Add the tool response(s) to history ---
            contents.append(types.Content(role="user", parts=tool_response_parts))
            print(f"Added {len(tool_response_parts)} tool response parts to history.")

            # --- 4.3 Make the next call to the model with updated history ---
            print("Making subsequent API call with tool responses...")
            response = await self.gemini.aio.models.generate_content(
                model=model,
                contents=contents,  # Send updated history
                config=types.GenerateContentConfig(
                    temperature=1.0,
                    tools=[tools],
                ),  # Keep sending same config
            )
            contents.append(response.candidates[0].content)

        if turn_count >= max_tool_turns and response.function_calls:
            print(f"Maximum tool turns ({max_tool_turns}) reached. Exiting loop.")

        print("MCP tool calling loop finished. Returning final response.")
        # --- 5. Return Final Response ---
        return response.text

    async def chat_loop(self):
        """Run an interactive chat loop"""
        print("\nMCP Client Started!")
        print("Type your queries or 'quit' to exit.")

        while True:
            try:
                query = input("\nQuery: ").strip()

                if query.lower() == "quit":
                    break

                response = await self.process_query(query)
                print("\n" + response)

            except Exception as e:
                print(f"\nError: {str(e)}")

    async def cleanup(self):
        """Clean up resources"""
        await self.exit_stack.aclose()


async def main():
    if len(sys.argv) < 2:
        print("Usage: python client.py <path_to_server_script>")
        sys.exit(1)

    client = MCPClient()
    try:
        await client.connect_to_server(sys.argv[1])
        await client.chat_loop()
    finally:
        await client.cleanup()


if __name__ == "__main__":
    import sys

    asyncio.run(main())
