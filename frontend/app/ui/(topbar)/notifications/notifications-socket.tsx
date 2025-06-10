'use client';

import { useState, useEffect } from 'react';

export default function useNotificationsSocket() {
  const [messages, setMessages] = useState<string[]>([]);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    const eventSource = new EventSource('/api/ws/subscribe');

    eventSource.onopen = () => {
      setConnected(true);
      console.log('✅ SSE connected');
    };

    eventSource.onmessage = (event) => {
      setMessages((prev) => [...prev, event.data]);
    };

    eventSource.onerror = (err) => {
      console.error('SSE error:', err);
      setConnected(false);
      eventSource.close();
    };

    return () => {
      eventSource.close();
    };
  }, []);


  const send = async (msg: string) => {
    await fetch('/api/ws/send', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ msg }),
    });
  };

  return { messages, connected, send };
}
