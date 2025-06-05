import { NextResponse } from 'next/server';

/**
 * Get all reviews created by user with given ID
 */
export async function POST(
  request: Request,
  { params }: { params: Promise<{ reviewerId: string }> },
) {
  const { reviewerId } = await params;
  console.log('reviewer id', reviewerId);

  const body = await request.text();
  const response = await fetch(
    `${process.env.API_URL}/review/reviewer/${reviewerId}`,
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
