from dotenv import load_dotenv
from pydantic_ai import Agent
from pydantic_ai.capabilities import WebSearch
from pydantic_ai_harness import CodeMode

load_dotenv()

agent = Agent(
    "openrouter:cohere/north-mini-code:free",
    capabilities=[
        # Wraps every tool into a single run_code tool, sandboxed by Monty
        # (https://github.com/pydantic/monty -- pulled in by the [code-mode] extra).
        # The model writes Python that calls multiple tools with loops, conditionals,
        # asyncio.gather, and local filtering -- one model round-trip for N tool calls.
        CodeMode(),
        # Provider-adaptive web search; native=False routes through the local
        # DuckDuckGo fallback (the [duckduckgo] extra above) so CodeMode can batch
        # web searches alongside the HN calls in a single run_code.
        WebSearch(local=True, native=False),
    ],
)

result = agent.run_sync(
    "Across the top, best, and 'show HN' Hacker News feeds, find the most-discussed "
    "story with at least 100 points. Pull its comment thread, its submitter's profile, "
    "and any web coverage. Summarize what you find in one paragraph."
)
print(result.output)
"""
The most-discussed HN story across top/best/show clearing 100 points is "Vibe coding
and agentic engineering are getting closer than I'd like" by Simon Willison (748 points,
853 comments, on the Best feed), submitted by long-time HNer e12e. The piece argues
that the two modes Willison once kept mentally separate -- throwaway "vibe coding" and
disciplined "agentic engineering" -- are blurring, since agents like Claude Code now
reliably handle non-trivial tasks like "build a JSON API endpoint that runs a SQL query"
with tests and docs on the first pass. The HN thread is unusually substantive, with
commenters debating whether LLMs created or merely *exposed* sloppy engineering
practices and warning of a "normalization of deviance" as engineers stop reviewing diffs.
"""
