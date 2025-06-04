// /api/ws/subscribe/route.ts
import { NextRequest } from 'next/server'
import { startNotificationsSocket, subscribe, unsubscribeAll } from '@/app/lib/notifications-socket'

const encoder = new TextEncoder()

export async function GET(request: NextRequest) {

    startNotificationsSocket()

    const stream = new ReadableStream({
        start(controller) {
            console.log('SSE stream started successfully')

            const onMessage = (msg: string) => {
                console.log(`Sending SSE message: ${msg.substring(0, 50)}${msg.length > 50 ? '...' : ''}`)
                controller.enqueue(encoder.encode(`data: ${msg}\n\n`))
            }

            const unsubscribe = subscribe(onMessage)

            request.signal.addEventListener('abort', () => {
                console.log('SSE connection aborted by client')
                unsubscribe()
                controller.close()
            })
        },
        cancel() {
            console.log('SSE stream cancelled')
            unsubscribeAll()
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
