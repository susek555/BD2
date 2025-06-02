'use client';

import { MouseEvent, useState } from 'react';
import { AverageRatingCard } from '../../review/average-rating-card';
import ReviewModal from '../../review/review-modal';

interface ReviewButtonProps {
  sellerRating?: number;
  sellerId?: number;
}

export default function ReviewButton({
  sellerRating,
  sellerId,
}: ReviewButtonProps) {
  const [dialogOpen, setDialogOpen] = useState(false);

  const onClick = (e: MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation();
    e.preventDefault();
    setDialogOpen(true);
    console.log('Review click');
  };

  return (
    <>
      {sellerRating ? (
        <button
          onClick={onClick}
          className='focus:ring-opacity-50 rounded-md transition-transform duration-200 hover:scale-105 focus:ring-2 focus:ring-blue-500 focus:outline-none'
        >
          <AverageRatingCard rating={sellerRating} />
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
        revieweeId={sellerId}
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        onSubmit={() => setDialogOpen(false)}
      />
    </>
  );
}
