import { Injectable } from '@angular/core';
import * as signalR from '@microsoft/signalr';
import { HttpClient } from '@angular/common/http';
import { MessageDTO } from '../DTO/MessageDTO';
import { Observable, Subject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ChatService {
  private connection: any = new signalR.HubConnectionBuilder()
    .withUrl('https://localhost:44375/chatsocket',) // fixed again
    .configureLogging(signalR.LogLevel.Information)
    .build();

  readonly POST_URL = 'https://localhost:7079/api/chat/send';

  private receivedMessageObject: MessageDTO = new MessageDTO();
  private sharedObj = new Subject<MessageDTO>();

  constructor(private http: HttpClient) {
    // invokes the function passed in when the connection with the server closes.
    this.connection.onclose(async () => {
      await this.start();
    });
    this.connection.on('ReceiveOne', (user: string, message: string) => {
      this.mapReceivedMessage(user, message);
    });
    this.start();
  }

  // Start the connection
  public async start() {
    try {
      await this.connection.start();
      console.log('connected');
    } catch (err) {
      console.log(err);
      setTimeout(() => this.start(), 5000);
    }
  }

  private mapReceivedMessage(user: string, message: string): void {
    this.receivedMessageObject.user = user;
    this.receivedMessageObject.msgText = message;
    this.sharedObj.next(this.receivedMessageObject);
  }

  // call the backend API through HTTP request that contains
  // the message sent from this iser
  public broadcastMessage(msgDTO: any) {
    this.http
      .post(this.POST_URL, msgDTO)
      .subscribe((data) => console.log(data));
  }

  // this method is used to share the data received from the backend
  // with other components in the project
  public retrieveMappedObject(): Observable<MessageDTO> {
    return this.sharedObj.asObservable();
  }
}
