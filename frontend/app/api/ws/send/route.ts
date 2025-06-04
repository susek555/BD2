import { NextRequest, NextResponse } from 'next/server'
import { sendNotification } from '@/app/lib/notifications-socket'

export async function POST(request: NextRequest) {
  try {
    const { msg } = await request.json()

    if (!msg || typeof msg !== 'string') {
      return NextResponse.json({ error: 'Invalid message' }, { status: 400 })
    }

    sendNotification(msg)

    return NextResponse.json({ success: true })
  } catch (error) {
    return NextResponse.json({ error: 'Failed to send message' }, { status: 500 })
  }
}
