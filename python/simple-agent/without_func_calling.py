from typing import List, Callable
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

    def get_tool_descriptions(self) -> str:
        return "\n".join([f"- {tool.name}: {tool.description}" for tool in self.tools])

    async def run(self, query: str) -> str:
        result = ""
        steps = []
        while True:
            prompt = f"""You are agent helping an answer a user query.
            Available tools: {self.get_tool_descriptions()}
G
            Original query: {query}
            Current results: {result}
            Steps taken: {", ".join(steps) if steps else "None"}

            If you have enough information to answer the query, response with: DONE: [final answer]
            Otherwise, respond with the name of the next tool to use (just the tool name).

            If you have no information to answer the query, and neither of the tools can help, respond with: DONE: [final answer]

            If you find yourself looping: respond with: DONE: [final answer]
            """

            print(f"prompt: {prompt}")
            response = await self.client.aio.models.generate_content(
                model="gemini-2.0-flash",
                contents=[types.Content(
                    role="user", parts=[types.Part(text=prompt)])],
                config=types.GenerateContentConfig(
                    max_output_tokens=100,
                    temperature=0,
                ),
            )

            answer = response.text.strip()
            print(answer)
            if answer.startswith("DONE:"):
                return answer[5:].strip()  # Everything after DONE:

            # Find and execute the matching tool
            tool_name = answer.lower()
            for tool in self.tools:
                if tool.name.lower() == tool_name:
                    tool_result = tool.func(query)
                    steps.append(f"Used {tool.name} (Step {len(steps) + 1})")
                    result += f"\nStep {len(steps)}: {tool_result}"
                    break


# Simulate tools (not actually implemented)
def check_calendar(task: str) -> str:
    return "Calendar shows: Meeting at 2pm"


def search_email(task: str) -> str:
    return "Found email from Bob about project deadline"


# Create tools list
tools = [
    Tool("calendar", "Checks calendar events", check_calendar),
    Tool("email", "Searches emails", search_email),
]


async def main():
    agent = Agent(tools)
    while True:
        query = input("What would you like to know? (or 'quit' to exit): ")
        if query.lower() == "quit":
            break

        result = await agent.run(query)
        print(result)


if __name__ == "__main__":
    asyncio.run(main())
