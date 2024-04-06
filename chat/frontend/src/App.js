import './App.css';
import { connect, sendMsg } from './api';
import React, { Component, useEffect } from "react";
import Header from './components/Header/Header'
import ChatHistory from './components/ChatHistory/ChatHistory'
import ChatInput from './components/ChatInput'
import Message from './components/Message/Message'

class App extends Component {
  constructor(props) {
    super(props);
    this.state = { chatHistory: [] };
  };

  componentDidMount() {
    connect((msg) => {
      console.log("New Message")
      this.setState(prevState => ({
        chatHistory: [...this.state.chatHistory, msg]
      }))
      console.log(this.state);
    });
  };

  send(event) {
    if (event.keyCode === 13) {
      sendMsg(event.target.value);
      event.target.value = "";
    }
  };

  render() {
    //这里不会改，应该是让chatHistory里的每个元素都封装成message，但是结果message是空
    const messages = this.state.chatHistory.map((msg) => {
      return <Message message={msg.data} />
    });

    return (
      <div className='ChatHistory'>
        <h2>Chat History</h2>
        {messages}
        <ChatInput send={this.send} />
      </div>
    );
  };
}
export default App;
