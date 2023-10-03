using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.SignalR;
using PublicChat.Hubs;
using PublicChat.ReqDTO;

namespace PublicChat.Controllers;

[Route("api/chat")]
[ApiController]
public class ChatController: ControllerBase
{
    private readonly IHubContext<ChatHub> _hubContext;

    public ChatController(IHubContext<ChatHub> hubContext)
    {
        _hubContext = hubContext;
    }

    [Route("send")]
    [HttpPost]
    public IActionResult SendRequest([FromBody] MessageDTO msg)
    {
        // IHubContext interface
        // The hub sends the message to all the Clients and invokes a client function
        // with a name similar to the one defined here as the first parameter
        // ( “ReceiveOne” ) in SendAsync(….) method.
        _hubContext.Clients.All.SendAsync("ReceiveOne", msg.user, msg.msgText);
        return Ok();
    }
}