import { fetchWithRefresh } from '@/app/lib/api/fetchWithRefresh';
import { NextRequest, NextResponse } from 'next/server';
import snakecaseKeys from 'snakecase-keys';

export async function PUT(request: NextRequest) {
  const userProfile = snakecaseKeys(await request.json());
  console.log(userProfile);

  const response = await fetchWithRefresh(`${process.env.API_URL}/users`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userProfile),
  });

  return new NextResponse(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers: response.headers,
  });
}
