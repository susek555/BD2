import { NextResponse } from 'next/server';

/**
 * Get all reviews created by user with given ID
 */
export async function POST(
  request: Request,
  { params }: { params: Promise<{ reviewerId: string }> },
) {
  const { reviewerId } = await params;

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
  // return NextResponse.json({
  //   reviews: [
  //     {
  //       id: 1,
  //       rating: 4,
  //       description: 'good',
  //       date: '2025-3-15',
  //       reviewer: { id: 1, username: 'damian' },
  //       reviewee: { id: 2, username: 'ggg' },
  //     },
  //   ],
  //   pagination: {
  //     total_pages: 1,
  //     total_records: 1,
  //   },
  // });
}
