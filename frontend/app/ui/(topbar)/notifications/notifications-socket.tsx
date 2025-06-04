'use client';

import { useState, useEffect, useRef } from 'react';

export default function useNotificationsSocket() {
  const socketRef = useRef<WebSocket | null>(null);
  const [messages, setMessages] = useState<string[]>([]);
  const [connected, setConnected] = useState(false);

  const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InN1c2Vsb3dza2lAbzIucGwiLCJ1c2VyX2lkIjoxLCJpc3MiOiJjYXItZGVhbGVyLWFwaSIsInN1YiI6IjEiLCJleHAiOjE3NTE2NDU1NzYsIm5iZiI6MTc0OTA1MzU3NiwiaWF0IjoxNzQ5MDUzNTc2fQ.e79ejUX49kBJStugJdhBx920wrlqbWEqGKlKyX2Balk"; // Replace with your actual token logic if needed

  useEffect(() => {


    let socket: WebSocket;
    try {
      socket = new WebSocket("ws://localhost:8080/ws", [token]);
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error);
      return;
    }
    socketRef.current = socket;

    socket.onopen = () => {
      setConnected(true);
      console.log('Socket connected');
    };

    socket.onmessage = (event) => {
      setMessages((prev) => [...prev, event.data]);
      console.log('Message received:', event.data);
    };

    socket.onclose = () => {
      setConnected(false);
      console.log('Socket disconnected');
    };

    return () => {
      socket.close();
    };
  }, []);

  const send = (msg: string) => {
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.send(msg);
    } else {
      console.warn('WebSocket not connected.');
    }
  };

  return { messages, connected, send };
}