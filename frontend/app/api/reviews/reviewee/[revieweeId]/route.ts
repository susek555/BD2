import { Review } from '@/app/lib/definitions/reviews';
import { NextResponse } from 'next/server';

/**
 * Get all reviews about user with given ID
 */
export async function GET(
  request: Request,
  { params }: { params: Promise<{ revieweeId: string }> },
) {
  const { revieweeId } = await params;
  const data: Review[] = [
    {
      id: 1,
      description:
        'Great service! Very responsive and professional. Would definitely recommend for any project.',
      rating: 5,
      date: '2023-09-15',
      reviewer: {
        id: 101,
        username: 'john_dev',
      },
      reviewee: {
        id: parseInt(revieweeId),
        username: 'techmaster',
      },
    },
    {
      id: 2,
      description:
        'Good work but took longer than expected. Quality was there though.',
      rating: 4,
      date: '2023-10-20',
      reviewer: {
        id: 102,
        username: 'sara_coder',
      },
      reviewee: {
        id: parseInt(revieweeId),
        username: 'techmaster',
      },
    },
    {
      id: 3,
      description: 'Average experience. Communication could have been better.',
      rating: 3,
      date: '2023-11-05',
      reviewer: {
        id: 103,
        username: 'dev_mike',
      },
      reviewee: {
        id: parseInt(revieweeId),
        username: 'techmaster',
      },
    },
    {
      id: 4,
      description:
        'Excellent problem-solving skills! Fixed a bug that had been bothering us for months.',
      rating: 5,
      date: '2023-12-10',
      reviewer: {
        id: 104,
        username: 'emma_tech',
      },
      reviewee: {
        id: parseInt(revieweeId),
        username: 'techmaster',
      },
    },
  ];
  return NextResponse.json(data);
}
