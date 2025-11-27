import { useState, useEffect, useRef, useCallback } from 'react';
import type { Topic } from '../types/topics';

const env = import.meta.env.VITE_APP_ENV;
const WS_HOST = env == "production" ? `ws://${document.location.hostname}:${document.location.port}` : "ws://localhost:8080";

export const useTopic = <T> (topic: Topic, onMessage: (msg: T) => void) => {
  const [isConnected, setIsConnected] = useState(false);
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    const ws = new WebSocket(`${WS_HOST}/ws`);
    socketRef.current = ws;

    ws.onopen = () => {
      console.log(`Connecté au WS, abonnement au topic : ${topic}`);
      setIsConnected(true);

      const payload = {
        action: "subscribe",
        topic: topic
      };
      ws.send(JSON.stringify(payload));
    };


    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);

        if (data.topic === topic) {
          onMessage(data.message);
        }
      } catch (err) {
        console.error("Error parsing JSON", err);
      }
    };

    ws.onclose = () => {
      console.log("disconected");
      setIsConnected(false);
    };

    return () => {
      ws.close();
    };
  }, [topic]); 

  const sendMessage = useCallback((msgContent: any) => {
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      const payload = {
        action: "publish",
        topic: topic,
        message: msgContent 
      };
      socketRef.current.send(JSON.stringify(payload));
    }
  }, [topic]);

  return { isConnected, sendMessage };
};