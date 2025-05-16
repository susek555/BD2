'use client';

import { Review } from '@/app/lib/definitions';
import { PencilIcon, StarIcon } from '@heroicons/react/20/solid';

interface ReviewCardProps {
  review: Review;
  variant: 'for' | 'by'; // 'for' = reviews for this user, 'by' = reviews by this user
}

export function ReviewCard({ review, variant }: ReviewCardProps) {
  const headerText =
    variant === 'for'
      ? `From ${review.reviewer.username}`
      : `Of ${review.reviewee.username}`;

  const handleEditClick = () => {
    console.log('Edit review:', review.id);
    alert(`Editing review for ${review.reviewee.username}`);
  };

  return (
    <div className='rounded-lg border border-gray-200 p-4 shadow-sm transition-shadow hover:shadow-md'>
      <div className='mb-2 flex items-center justify-between'>
        <span className='text-sm font-medium text-gray-700'>{headerText}</span>
        <div className='flex items-center gap-2'>
          <span className='text-xs text-gray-500'>{review.date}</span>
          {variant === 'by' && (
            <button
              onClick={handleEditClick}
              className='mb-1 ml-1 rounded-full p-1 transition-colors hover:bg-gray-100'
              aria-label='Edit review'
            >
              <PencilIcon className='h-4 w-4 text-gray-500' />
            </button>
          )}
        </div>
      </div>
      <div className='mb-2'>
        <div className='flex items-center'>
          {Array.from({ length: 5 }).map((_, i) => (
            <StarIcon
              key={i}
              className={`h-4 w-4 ${i < review.rating ? 'text-yellow-400' : 'text-gray-300'}`}
            />
          ))}
          <span className='mt-1 ml-2 text-xs text-gray-500'>
            ({review.rating}/5)
          </span>
        </div>
      </div>
      <p className='text-sm text-gray-600'>{review.description}</p>
    </div>
  );
}
