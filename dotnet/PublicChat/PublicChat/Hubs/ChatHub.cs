using Microsoft.AspNetCore.SignalR;

namespace PublicChat.Hubs;

public class ChatHub: Hub
{
    public Task SendMessage(string user, string message)
    {
        // send the user and message to the client-side receiver with the name "ReceiveOne"
        return Clients.All.SendAsync("ReceiveOne", user, message);
    }
}