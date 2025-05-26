import { NextResponse } from 'next/server';

export async function GET() {
  return NextResponse.json({
    status: 'ok',
    message: 'Frontend is healthy',
    timestamp: new Date().toISOString()
  });
}
