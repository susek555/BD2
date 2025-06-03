import { NextResponse } from 'next/server';

/**
 * Get review created about reviewee with given ID by reviewer with given ID
 */
export async function GET(
  request: Request,
  { params }: { params: Promise<{ revieweeId: string; reviewerId: string }> },
) {
  const { reviewerId, revieweeId } = await params;
  const review = {
    id: 1,
    description:
      'Great service! Very responsive and professional. Would definitely recommend for any project.',
    rating: 5,
    date: '2023-09-15',
    reviewer: {
      id: parseInt(reviewerId),
      username: 'john_dev',
    },
    reviewee: {
      id: parseInt(revieweeId),
      username: 'techmaster',
    },
  };
  
  return NextResponse.json(review);
}
