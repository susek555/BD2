import { fetchWithRefresh } from '@/app/lib/api/fetchWithRefresh';
import { NextRequest, NextResponse } from 'next/server';

export async function PUT(request: NextRequest) {
  const passwordData = await request.json();
  const response = await fetchWithRefresh(
    `${process.env.API_URL}/auth/change-password`,
    {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(passwordData),
    },
  );

  return new NextResponse(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers: response.headers,
  });
}
