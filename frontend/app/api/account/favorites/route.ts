import { fetchWithRefresh } from '@/app/lib/api/fetchWithRefresh';
import { NextResponse } from 'next/server';

export async function POST(request: Request) {
  const body = await request.text();
  const response = await fetchWithRefresh(
    `${process.env.API_URL}/sale-offer/filtered`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: body,
    },
  );

  const data = await response.json();
  return NextResponse.json(data);
}
