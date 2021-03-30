import graphene
from graphene_django import DjangoObjectType
from graphql import GraphQLError

from .models import Link, Vote
from users.schema import UserType


class LinkType(DjangoObjectType):
    class Meta:
        model = Link

class VoteType(DjangoObjectType):
    class Meta:
        model = Vote

class Query(graphene.ObjectType):
    links = graphene.List(LinkType,
                          search=graphene.String(),
                          first=graphene.Int(),
                          skip=graphene.Int())
    votes = graphene.List(VoteType)

    def resolve_links(self, info, search=None, first=None,
                      skip=None, **kwargs):
        qs = Link.objects.all()
        if search:
            filter = (
                Q(url__icontains=search) |
                Q(description__icontains=search)
            )
            qs = qs.filter(filter)

        if skip:
            qs = qs[skip:]
        if first:
            qs = qs[:first]
        return qs

    def resolve_votes(self, info, **kwargs):
        return Vote.objects.all()

# A mutation class
class CreateLink(graphene.Mutation):
    id = graphene.Int()
    url = graphene.String()
    description = graphene.String()
    posted_by = graphene.Field(UserType)

    # The data you can send to the server
    class Arguments:
        url = graphene.String()
        description = graphene.String()

    # Mutation method: create a link in the database using the data
    # sent by the user.
    def mutate(self, info, url, description):
        user = info.contex.user or None
        link = Link(url=url, description=description, posted_by=user)
        link.save()

        return CreateLink(
            id=link.id,
            url=link.url,
            description=link.description,
            posted_by=link.posted_by,
        )

class CreateVote(graphene.Mutation):
    user = graphene.Field(UserType)
    link = graphene.Field(LinkType)

    class Arguments:
        link_id = graphene.Int()

    def mutate(self, info, link_id):
        user = info.context.user
        if user.is_anonymous:
            raise GraphQLError('You must be logged to vote')

        link = Link.objects.filter(id=link_id).first()
        if not link:
            raise Exception('Invalid link!')

        Vote.objects.create(
            user=user,
            link=link
        )
        return CreateVote(user=user, link=link)

# Create a mutation class with a field to be resolved, which points
# to our mutation defined before.
class Mutation(graphene.ObjectType):
    create_link = CreateLink.Field()
    create_vote = CreateVote.Field()
