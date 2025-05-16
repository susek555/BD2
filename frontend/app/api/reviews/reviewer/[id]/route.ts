import { Review } from '@/app/lib/definitions';
import { NextResponse } from 'next/server';

/**
 * Get all reviews created by user with given ID
 */
export async function GET(
  request: Request,
  { params }: { params: Promise<{ id: string }> },
) {
  const { id } = await params;
  const data: Review[] = [
    {
      id: 1,
      description:
        'Great service! Very responsive and professional. Would definitely recommend for any project.',
      rating: 5,
      date: '2023-09-15',
      reviewer: {
        id: parseInt(id),
        username: 'john_dev',
      },
      reviewee: {
        id: 201,
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
        id: parseInt(id),
        username: 'john_dev',
      },
      reviewee: {
        id: 202,
        username: 'codewizard',
      },
    },
    {
      id: 3,
      description: 'Average experience. Communication could have been better.',
      rating: 3,
      date: '2023-11-05',
      reviewer: {
        id: parseInt(id),
        username: 'john_dev',
      },
      reviewee: {
        id: 203,
        username: 'devguru',
      },
    },
    {
      id: 4,
      description:
        'Excellent problem-solving skills! Fixed a bug that had been bothering us for months.',
      rating: 5,
      date: '2023-12-10',
      reviewer: {
        id: parseInt(id),
        username: 'john_dev',
      },
      reviewee: {
        id: 204,
        username: 'webmaster',
      },
    },
  ];
  return NextResponse.json(data);
}
