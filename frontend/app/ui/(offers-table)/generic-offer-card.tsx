import { BaseOffer } from '@/app/lib/definitions/SaleOffer';
import { CarImageHomePageSkeleton } from '@/app/ui/skeletons';
import { Suspense } from 'react';
import MainImage from './miniature_image';

export type OfferVariant = 'default' | 'history';

export interface GenericOfferProps<T extends BaseOffer = BaseOffer> {
  offer: T;
  variant?: OfferVariant;
  priceLabel?: string;
  headerContent?: React.ReactNode;
  className?: string;
}

export default function GenericOfferCard<T extends BaseOffer = BaseOffer>({
  offer,
  variant = 'default',
  priceLabel,
  headerContent,
  className = '',
}: GenericOfferProps<T>) {
  const getPriceLabel = () => {
    if (priceLabel) return priceLabel;

    switch (variant) {
      case 'history':
        return offer.is_auction ? 'Final bid' : 'Price';
      default:
        return offer.is_auction ? 'Current bid' : 'Price';
    }
  };

  return (
    <div className={`flex flex-row gap-4 ${className}`}>
      <div className='p-2.5'>
        <Suspense fallback={<CarImageHomePageSkeleton />}>
          <MainImage url={offer.main_url} />
        </Suspense>
      </div>

      <div className='flex h-40 w-full flex-col py-2.5'>
        <div className='mb-3 flex items-center justify-between'>
          <div className='flex items-center gap-2'>
            <h3 className='text-2xl font-bold'>{offer.name}</h3>
            {offer.is_auction && (
              <span className='mt-[-8] rounded-full bg-blue-100 px-3 py-1 text-sm font-semibold text-blue-800'>
                Auction
              </span>
            )}
          </div>
          {headerContent}
        </div>

        <div className='flex flex-1 flex-row'>
          <div className='flex-1 px-4'>
            <p className='text-bg mb-2'>
              Production year:{' '}
              <span className='font-bold'>
                {offer.production_year.toString()}
              </span>
            </p>
            <p className='text-bg'>
              Color: <span className='font-bold'>{offer.color}</span>
            </p>
          </div>

          <div className='flex-1'>
            <p className='text-bg'>
              Mileage:{' '}
              <span className='font-bold'>
                {offer.mileage.toLocaleString()} km
              </span>
            </p>
          </div>

          <div className='flex items-end pr-4'>
            <div className='text-right'>
              <p className='text-sm text-gray-600'>{getPriceLabel()}</p>
              <p className='text-2xl font-bold text-green-600'>
                {offer.price.toString()} zł
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
