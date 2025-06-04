// /api/ws/subscribe/route.ts
import { NextRequest } from 'next/server'
import { getServerSession } from 'next-auth'
import { authConfig } from '@/app/lib/authConfig'
import WebSocket from 'ws'

const encoder = new TextEncoder()

export async function GET(request: NextRequest) {
  const session = await getServerSession(authConfig)
  const token = session?.user?.refreshToken

  if (!token) {
    return new Response('Unauthorized', { status: 401 })
  }

  let socket: WebSocket | null = null

  const stream = new ReadableStream({
    async start(controller) {
      socket = new WebSocket('ws://localhost:8080/ws', [token])

      socket.on('open', () => {
        console.log('✅ WS connected for user:', session?.user?.email)
      })

      socket.on('message', (data) => {
        const msg = data.toString()
        controller.enqueue(encoder.encode(`data: ${msg}\n\n`))
      })

      socket.on('close', () => {
        console.log('❌ WS closed for user:', session?.user?.email)
        controller.close()
      })

      socket.on('error', (err) => {
        console.error('WebSocket error:', err)
        controller.error(err)
      })

      request.signal.addEventListener('abort', () => {
        if (socket) {
          socket.close()
        }
      })
    },
    cancel() {
      if (socket) {
        socket.close()
      }
    },
  })

  return new Response(stream, {
    headers: {
      'Content-Type': 'text/event-stream',
      'Cache-Control': 'no-cache',
      Connection: 'keep-alive',
    },
  })
}

