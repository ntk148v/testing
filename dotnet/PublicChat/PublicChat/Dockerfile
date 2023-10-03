FROM mcr.microsoft.com/dotnet/aspnet:6.0 AS base
WORKDIR /app
EXPOSE 80
EXPOSE 443

FROM mcr.microsoft.com/dotnet/sdk:6.0 AS build
WORKDIR /src
COPY ["PublicChat/PublicChat.csproj", "PublicChat/"]
RUN dotnet restore "PublicChat/PublicChat.csproj"
COPY . .
WORKDIR "/src/PublicChat"
RUN dotnet build "PublicChat.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "PublicChat.csproj" -c Release -o /app/publish /p:UseAppHost=false

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "PublicChat.dll"]
