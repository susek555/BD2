import { NextRequest, NextResponse } from 'next/server'
import { NOTIFICATIONS_SOCKET_URL } from '@/app/lib/constants'

export async function POST(request: NextRequest) {
  try {
    const { msg } = await request.json()

    if (!msg || typeof msg !== 'string') {
      return NextResponse.json({ error: 'Invalid message' }, { status: 400 })
    }

    const response = await fetch(`${NOTIFICATIONS_SOCKET_URL}/send`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ msg }),
    });

    return response.json()
  } catch (error) {
    return NextResponse.json({ error: 'Failed to send message' }, { status: 444 })
  }
}
