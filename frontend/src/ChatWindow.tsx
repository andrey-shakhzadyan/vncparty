function ChatWindow() {
	return (
	  <div className="basis-20 min-h-full flex-col min-w-1/6 shadow-sm resize-x">
	    <div className="bg-100 hover:display-none rounded-sm p-1 min-w-full">
	      <div className="">
		      Chat
	      </div>
	    </div>
	      <div className="flex-col focus:outline-2 min-h-4/5 z-0">
		<div className="chat chat-start">
		  <div className="chat-header">
						 Chatter A
		    <time className="text-xs opacity-50">2 hours ago</time>
		  </div>
		  <div className="chat-bubble">Test message 1</div>
		  <div className="chat-footer opacity-50">Seen</div>
		</div>
		<div className="chat chat-end">
		  <div className="chat-header">
						 Chatter B
		    <time className="text-xs opacity-50">2 hour ago</time>
		  </div>
		  <div className="chat-bubble">I disagree.</div>
		  <div className="chat-footer opacity-50">Delivered</div>
		</div>
	      </div>
	    <input type="text" placeholder="..." className="flex-auto input outline-none" />
	  </div>
	)
}

export default ChatWindow;
