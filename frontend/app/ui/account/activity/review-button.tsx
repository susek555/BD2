'use client';

import { Review } from '@/app/lib/definitions/reviews';
import { MouseEvent, useState } from 'react';
import { AverageRatingCard } from '../../review/average-rating-card';
import ReviewModal from '../../review/review-modal';

interface ReviewButtonProps {
  sellerId: number;
  userId: number;
  review?: Review | null;
}

export default function ReviewButton({
  sellerId,
  userId,
  review,
}: ReviewButtonProps) {
  const [dialogOpen, setDialogOpen] = useState(false);

  const onClick = async (e: MouseEvent<HTMLButtonElement>) => {
    console.log('Review click');
    e.stopPropagation();
    e.preventDefault();
    setDialogOpen(true);
  };

  return (
    <>
      {review && !('error_description' in review) ? (
        <button
          onClick={onClick}
          className='focus:ring-opacity-50 rounded-md transition-transform duration-200 hover:scale-105 focus:ring-2 focus:ring-blue-500 focus:outline-none'
        >
          <AverageRatingCard rating={review.rating} />
        </button>
      ) : (
        <button
          onClick={onClick}
          className='focus:ring-opacity-50 rounded-md bg-blue-600 px-4 py-2 text-white transition-colors duration-200 hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:outline-none'
        >
          Rate Seller
        </button>
      )}
      <ReviewModal
        review={review}
        revieweeId={sellerId}
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        onSubmit={() => {
          setDialogOpen(false);
          window.location.reload();
        }}
      />
    </>
  );
}
