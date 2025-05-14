import { SaleOffer } from '@/app/lib/definitions';
import { CarImageSkeleton } from '@/app/ui/skeletons';
import { Suspense } from 'react';

export default function SingleOffer({ offer }: { offer: SaleOffer }) {
  return (
    <div className='flex flex-row gap-4'>
      <div className='p-2.5'>
        <Suspense fallback={<CarImageSkeleton />}>
          <CarImageSkeleton />
          {/* TODO implement fetching image from internet */}
        </Suspense>
      </div>

      <div className='flex h-40 w-full flex-col py-2.5'>
        <div className='mb-3 flex items-center justify-between'>
          <div className='flex items-center gap-2'>
            <h3 className='text-2xl font-bold'>{offer.name}</h3>
            {offer.isAuction && (
              <span className='mt-[-8] rounded-full bg-blue-100 px-3 py-1 text-sm font-semibold text-blue-800'>
                Auction
              </span>
            )}
          </div>
        </div>

        <div className='flex flex-1 flex-row'>
          <div className='flex-1 px-4'>
            <p className='text-bg mb-2'>
              Production year:{' '}
              <span className='font-bold'>
                {offer.productionYear.toString()}
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

          <div className='flex h-full items-end justify-end self-end'>
            <p className='px-3 text-2xl font-bold'>
              {offer.price.toString()} zł
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
