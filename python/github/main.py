import os

from github import Github, InputGitAuthor

# First create a Github instance:
# using an access token
g = Github(os.environ.get("GITHUB_TOKEN"))

me = g.get_user()

print(g.get_user().add_to_starred(g.get_repo("kavu/go_reuseport")))

repo = g.get_repo("ntk148v/til")
author = InputGitAuthor("Kien Nguyen Tuan", "kiennt2609@gmail.com")
with open('/home/kiennt65/Workspace/github.com/ntk148v/til/nomad/README.md', 'r') as f:
    data = f.read()
    print(repo.create_file("nomad/README.md", "Add Nomad til", data, author=author, committer=author))
