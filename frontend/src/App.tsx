import { useRef, useState, useEffect } from "react";

import { VncScreen } from "react-vnc";
import Header from "./Header.tsx";
import ChatWindow from "./ChatWindow.tsx";
import { FaLock } from "react-icons/fa6";

function App() {
  const ref = useRef<React.ElementRef<typeof VncScreen>>(null);

  const params = new URLSearchParams(window.location.search);
  const uuid = params.get("uuid");

  const [vncAddr, setVncAddr] = useState("");
  const [roomName, setRoomName] = useState("");

  const [screenLocked, setScreenLocked] = useState(false);

  useEffect(() => {
    fetch("/api/get_room?uuid=" + uuid)
      .then((response) => response.json())
      .then((roomData) => {
        setVncAddr("ws://localhost:1454/ws/roomproxy/" + uuid);
        setRoomName(roomData.room_name);
        if (ref.current) {
          ref.current.connect();
        }
      });
  });
  function VNCReconnect() {
    if (ref.current?.connected) {
      ref.current.disconnect();
    }
    ref.current?.connect();
  }

  function VNCShutdown() {
    if (ref.current?.connected) {
      ref.current.machineShutdown();
    }
  }

  function VNCReboot() {
    if (ref.current?.connected) {
      ref.current.machineReboot();
    }
  }

  function toggleScreenLock() {
    setScreenLocked(!screenLocked);
  }

  return (
    <div>
      <Header roomName={roomName} />
      <br />
      <div className="flex justify-items-center p-2">
        <div className="flex flex-col basis-60 h-dvh grow  items-center">
          <div>
            {/*{!ref.current?.connected && <span className="z-4 absolute loading loading-spinner loading-xl m-auto"></span>} */}
            <div>
              {screenLocked && (
                <div className="interactionblocker z-4 absolute"> </div>
              )}
              <div className="indicator">
                {screenLocked && (
                  <span className="indicator-item badge badge-primary z-4">
                    {" "}
                    <FaLock />{" "}
                  </span>
                )}
                <VncScreen
                  className="z-2 relative basis-90 grow loadingstatic shadow-sm"
                  url={vncAddr}
                  scaleViewport
                  background="white"
                  style={{
                    width: "121vh",
                    height: "85.25vh",
                  }}
                  ref={ref}
                />
              </div>
            </div>
          </div>
          <br />
          <div className="join">
            <button className="btn btn-error join-item" onClick={VNCShutdown}>
              Shut Down
            </button>
            <button className="btn join-item" onClick={VNCReboot}>
              Reboot
            </button>
            <button className="btn join-item" onClick={VNCReconnect}>
              Reconnect
            </button>
            <button className="btn join-item" onClick={toggleScreenLock}>
              Toggle screen lock
            </button>
          </div>
          <br />
        </div>
        <ChatWindow />
      </div>
    </div>
  );
}

export default App;
