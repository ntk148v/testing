import { Component, OnInit } from '@angular/core';
import { ChatService } from './services/chat.service';
import { MessageDTO } from './DTO/MessageDTO';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  constructor(private chatService: ChatService) {}

  ngOnInit(): void {
    this.chatService
      .retrieveMappedObject()
      .subscribe((receivedObj: MessageDTO) => {
        this.addToInbox(receivedObj);
      }); // calls the service method to get the new messages sent
  }

  msgDTO: MessageDTO = new MessageDTO();
  msgInboxArray: MessageDTO[] = [];

  send(): void {
    if (this.msgDTO) {
      if (this.msgDTO.user.length == 0 || this.msgDTO.user.length == 0) {
        window.alert('Both fields are required.');
        return;
      } else {
        this.chatService.broadcastMessage(this.msgDTO); // Send the message via a service
        this.msgDTO.msgText = '';
      }
    }
  }

  addToInbox(obj: MessageDTO) {
    let newObj = new MessageDTO();
    newObj.user = obj.user;
    newObj.msgText = obj.msgText;
    this.msgInboxArray.push(newObj);
  }
}
