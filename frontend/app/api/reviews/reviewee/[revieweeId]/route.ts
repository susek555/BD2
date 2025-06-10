import { NextResponse } from 'next/server';

/**
 * Get all reviews about user with given ID
 */
export async function POST(
  request: Request,
  { params }: { params: Promise<{ revieweeId: string }> },
) {
  const { revieweeId } = await params;
  const body = await request.text();
  const response = await fetch(
    `${process.env.API_URL}/review/reviewee/${revieweeId}`,
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
