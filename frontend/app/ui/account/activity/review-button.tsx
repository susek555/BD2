'use client';

import { getReviewByRevieweeReviewer } from '@/app/lib/api/reviews';
import { Review } from '@/app/lib/definitions/reviews';
import { MouseEvent, useState } from 'react';
import { AverageRatingCard } from '../../review/average-rating-card';
import ReviewModal from '../../review/review-modal';
import toast from 'react-hot-toast';

interface ReviewButtonProps {
  sellerRating?: number;
  sellerId: number;
  userId: number;
}

export default function ReviewButton({
  sellerRating,
  sellerId,
  userId,
}: ReviewButtonProps) {
  const [dialogOpen, setDialogOpen] = useState(false);
  const [review, setReview] = useState<Review | undefined>(undefined);

  const onClick = async (e: MouseEvent<HTMLButtonElement>) => {
    console.log('Review click');
    e.stopPropagation();
    e.preventDefault();
    try {
      if (sellerRating) {
        const review = await getReviewByRevieweeReviewer(sellerId, userId);
        setReview(review);
      }
    } catch {
      toast.error("Something went wrong fetching review")
      return;
    }
    setDialogOpen(true);
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
