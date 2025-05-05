from typing import List, Dict, Callable
from dataclasses import dataclass

import asyncio
import dotenv

from google import genai
from google.genai import types

dotenv.load_dotenv()


@dataclass
class Tool:
    name: str
    description: str
    func: Callable


class Agent:
    def __init__(self, tools: List[Tool]):
        self.tools = tools
        self.client = genai.Client()
        self.model = "gemini-2.0-flash"

    def get_tools(self) -> types.Tool:
        return types.Tool(
            function_declarations=[
                {
                    "name": "calendar",
                    "description": "Checks calendar events for a person",
                    "parameters": {
                        "type": "object",
                        "properties": {
                            "person_name": {
                                "type": "string",
                                "description": "The name of the person to look up"
                            }
                        },
                        "required": ["person_name"]
                    }
                },
                {
                    "name": "email",
                    "description": "Searches emails for a person",
                    "parameters": {
                        "type": "object",
                        "properties": {
                            "person_name": {
                                "type": "string",
                                "description": "The name of the person to look up"
                            }
                        },
                        "required": ["person_name"]
                    }
                }
            ]
        )

    async def run(self, query: str) -> str:
        result = ""
        steps = []

        # Initialize chat history
        chat_history = []

        # Add initial message to history
        initial_message = types.Content(
            role="user",
            parts=[types.Part(text=f"Help answer this query: {query}\n"
                              f"Respond with DONE: [answer] if no tools needed, we're looping, "
                              f"or no tools left to call for new info. Be thorough and provide "
                              f"all answers you possibly can if there are multiple tools.")]
        )

        chat_history.append(initial_message)

        while True:
            print("looping")
            # Generate content
            response = await self.client.aio.models.generate_content(
                model=self.model,
                contents=chat_history,
                config=types.GenerateContentConfig(
                    temperature=0,
                    max_output_tokens=100,
                    tools=[self.get_tools()],
                )
            )

            # Get the response text
            message_text = response.text

            # Get the response text
            message_text = response.text

 # Check if we're done
            if message_text and message_text.startswith("DONE:"):
                return message_text[5:].strip()

            # Create assistant message for history
            assistant_message = types.Content(
                role="assistant",
                parts=[types.Part(text=message_text)]
            )

            # Check for function calls
            function_called = False

            if hasattr(response, "candidates") and response.candidates:
                for candidate in response.candidates:
                    if hasattr(candidate, "content") and candidate.content:
                        for part in candidate.content.parts:
                            if hasattr(part, "function_call"):
                                function_call = part.function_call
                                if not function_call:
                                    continue
                                tool_name = function_call.name

                                # Extract person_name from args
                                args = function_call.args
                                person_name = None

                                # Handle different ways args might be structured
                                if isinstance(args, dict):
                                    person_name = args.get("person_name")
                                elif hasattr(args, "get"):
                                    person_name = args.get("person_name")

                                # Skip if no person_name provided
                                if not person_name:
                                    continue

                                # Find and execute the matching tool
                                for tool in self.tools:
                                    if tool.name == tool_name:
                                        tool_result = tool.func(person_name)
                                        steps.append(
                                            f"Used {tool.name} for {person_name}")
                                        result += f"\n{tool_result}"

                                        # Add assistant message to history (including function call)
                                        chat_history.append(assistant_message)

                                        # Add function response to history
                                        tool_response = types.Content(
                                            role="user",
                                            parts=[types.Part(text=f"I used {tool.name} for {person_name} and got this result: {tool_result}\n"
                                                              f"Continue or respond with DONE: [answer] if we have enough information.")]
                                        )

                                        chat_history.append(tool_response)
                                        function_called = True
                                        break

            # If no function was called, we might be done
            if not function_called:
                if message_text:
                    return message_text
                else:
                    return f"Could not process the query: {query}"


# Simulate tools with name-based responses
def check_calendar(name: str) -> str:
    calendar_data = {
        "alice": "Meeting with clients at 2pm",
        "bob": "Team standup at 10am",
        "charlie": "Lunch meeting at 12pm",
        "default": "No meetings found"
    }
    return f"Calendar for {name}: {calendar_data.get(name.lower(), calendar_data['default'])}"


def search_email(name: str) -> str:
    email_data = {
        "alice": "Latest email about Q4 planning",
        "bob": "Project status update",
        "charlie": "Vacation request pending",
        "default": "No recent emails"
    }
    return f"Emails from {name}: {email_data.get(name.lower(), email_data['default'])}"


# Create tools list
tools = [
    Tool("calendar", "Checks calendar events for a person", check_calendar),
    Tool("email", "Searches emails for a person", search_email),
]


async def main():
    agent = Agent(tools)

    while True:
        query = input("What would you like to know? (or 'quit' to exit): ")
        if query.lower() == 'quit':
            break

        result = await agent.run(query)
        print(result)


if __name__ == "__main__":
    asyncio.run(main())
