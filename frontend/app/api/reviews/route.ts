import { fetchWithRefresh } from '@/app/lib/api/fetchWithRefresh';
import { NextRequest } from 'next/server';

export async function POST(request: NextRequest) {
  const body = await request.text();
  const response = fetchWithRefresh(`${process.env.API_URL}/review`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: body,
  });
}

export async function PUT(request: NextRequest) {
  const body = await request.text();
  const response = fetchWithRefresh(`${process.env.API_URL}/review`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: body,
  });
}
