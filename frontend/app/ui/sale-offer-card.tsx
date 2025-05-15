'use client';

import { updateFavoriteStatus } from '@/app/lib/api/listing/requests';
import { SaleOffer } from '@/app/lib/definitions';
import GenericOfferCard, {
  GenericOfferProps,
} from '@/app/ui/generic-offer-card';
import { StarIcon as StarOutline } from '@heroicons/react/24/outline';
import { StarIcon as StarSolid } from '@heroicons/react/24/solid';
import { useState } from 'react';

export default function SaleOfferCard({
  offer,
  variant = 'default',
  priceLabel,
  className = '',
}: GenericOfferProps<SaleOffer>) {
  const [isFavorite, setIsFavorite] = useState(offer.isFavorite);

  const handleFavoriteToggle = () => {
    updateFavoriteStatus(offer.id, offer.isFavorite);
    setIsFavorite(!isFavorite);
    offer.isFavorite = !isFavorite;
  };

  const favoriteButton = (
    <button
      onClick={(e) => {
        e.stopPropagation();
        e.preventDefault();
        handleFavoriteToggle();
      }}
      className='mr-2 ml-auto rounded-full p-1.5 transition-colors hover:bg-gray-100'
      aria-label={
        offer.isFavorite ? 'Remove from favorites' : 'Add to favorites'
      }
    >
      {isFavorite ? (
        <StarSolid className='h-6 w-6 text-yellow-500' />
      ) : (
        <StarOutline className='h-6 w-6 text-gray-500 hover:text-yellow-400' />
      )}
    </button>
  );

  return (
    <GenericOfferCard
      offer={offer}
      variant={variant}
      priceLabel={priceLabel}
      headerContent={favoriteButton}
      className={className}
    />
  );
}
