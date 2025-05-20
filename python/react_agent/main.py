from langgraph.prebuilt import create_react_agent
from langchain_community.utilities import SerpAPIWrapper
from langchain_core.tools import Tool
from dotenv import load_dotenv

load_dotenv()
search = SerpAPIWrapper(params={"engine": "google"})
tools = [
    Tool(
        name="web_search",
        description="Search the web for information",
        func=search.run,
    )
]
agent = create_react_agent(model="google_genai:gemini-2.0-flash", tools=tools)
input_message = {
    "role": "user",
    "content": "How many kids do the band members of Metallica have?",
}

for step in agent.stream(
    {"messages": [input_message]},
    stream_mode="values",
):
    step["messages"][-1].pretty_print()
