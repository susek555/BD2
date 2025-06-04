// lib/notificationsSocket.ts
import WebSocket from 'ws';

let socket: WebSocket | null = null;
let subscribers: ((msg: string) => void)[] = [];

export function startNotificationsSocket() {
  if (socket) return;

  const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InN1c2Vsb3dza2lAbzIucGwiLCJ1c2VyX2lkIjoxLCJpc3MiOiJjYXItZGVhbGVyLWFwaSIsInN1YiI6IjEiLCJleHAiOjE3NTE2NTc1MTIsIm5iZiI6MTc0OTA2NTUxMiwiaWF0IjoxNzQ5MDY1NTEyfQ.P6R96oMqHVw5VpCpXo6e1YvWFMcDP-FyTRCmq3ApokA';
  socket = new WebSocket(`ws://localhost:8080/ws`, [token]);

  socket.on('open', () => {
    console.log('✅ WebSocket server-side connected');
  });

  socket.on('message', (data) => {
    const msg = data.toString();
    console.log('📨 Serwer otrzymał wiadomość:', msg);
    subscribers.forEach((fn) => fn(msg)); // powiadom subskrybentów
  });

  socket.on('close', () => {
    console.log('❌ WebSocket zamknięty');
    socket = null;
  });

  socket.on('error', (err) => {
    console.error('WebSocket error:', err);
  });
}

export function sendNotification(msg: string) {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(msg);
  }
}

export function subscribe(fn: (msg: string) => void) {
  subscribers.push(fn);
  return () => {
    subscribers = subscribers.filter((f) => f !== fn);
  };
}

export function unsubscribeAll() {
  subscribers = [];
}
