'use client';
import { StarIcon as Star } from '@heroicons/react/20/solid';
import { ReactElement } from 'react';

interface AverageRatingCardProps {
  rating: number;
  starSize?: number;
}

export function AverageRatingCard({
  rating,
  starSize = 20,
}: AverageRatingCardProps) {
  const renderStars = (): ReactElement[] => {
    const stars: ReactElement[] = [];
    const fullStars = Math.floor(rating);
    const partialStar = rating % 1;
    const emptyStars = 5 - fullStars - (partialStar > 0 ? 1 : 0);

    for (let i = 0; i < fullStars; i++) {
      stars.push(
        <Star
          key={`full-${i}`}
          className='text-yellow-500'
          width={starSize}
          height={starSize}
        />,
      );
    }

    if (partialStar > 0) {
      stars.push(
        <div key='partial' className='relative'>
          <Star className='text-gray-300' width={starSize} height={starSize} />
          <div
            className='absolute top-0 left-0 overflow-hidden'
            style={{ width: `${partialStar * 100}%` }}
          >
            <Star
              className='text-yellow-500'
              width={starSize}
              height={starSize}
            />
          </div>
        </div>,
      );
    }

    for (let i = 0; i < emptyStars; i++) {
      stars.push(
        <Star
          key={`empty-${i}`}
          className='text-gray-300'
          width={starSize}
          height={starSize}
        />,
      );
    }

    return stars;
  };

  return (
    <div className='flex items-center gap-2'>
      <div className='flex'>{renderStars()}</div>
      <span className='ml-2 text-lg font-semibold'>{rating.toFixed(1)}</span>
    </div>
  );
}
