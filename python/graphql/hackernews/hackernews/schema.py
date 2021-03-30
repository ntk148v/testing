import graphene
import graphql_jwt

import links.schema
import links.schema_relay
import users.schema


class Query(users.schema.Query,
            links.schema.Query,
            links.schema_relay.Query,
            graphene.ObjectType):
    pass

class Mutation(users.schema.Mutation, links.schema.Mutation, graphene.ObjectType):
    token_auth = graphql_jwt.ObtainJSONWebToken.Field()
    verify_token = graphql_jwt.Verify.Field()
    refresh_token = graphql_jwt.Verify.Field()

schema = graphene.Schema(query=Query, mutation=Mutation)
