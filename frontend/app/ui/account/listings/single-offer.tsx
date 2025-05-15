'use client';

import { deleteListing } from '@/app/lib/api/listing/requests';
import { SaleOffer } from '@/app/lib/definitions';
import { CarImageSkeleton } from '@/app/ui/skeletons';
import { PencilIcon, TrashIcon } from '@heroicons/react/20/solid';
import { Suspense, useState } from 'react';
import ConfirmationModal from '../../confirm-modal';

export default function SingleListingsOffer({ offer }: { offer: SaleOffer }) {
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  const handleDelete = () => {
    deleteListing(offer.id);
    setShowDeleteConfirm(false);
  };

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
          <div className='flex space-x-2 pr-4'>
            <button
              className='rounded-md bg-blue-600 px-3 py-1 text-sm font-medium text-white transition-colors hover:bg-blue-700'
              onClick={(e) => {
                e.stopPropagation();
                e.preventDefault();
                console.log('Edit offer:', offer.id); // Todo implement redirect to edit offer page
              }}
            >
              <PencilIcon className='h-5 w-5' />
            </button>
            <button
              className='rounded-md bg-red-600 px-3 py-1 text-sm font-medium text-white transition-colors hover:bg-red-700'
              onClick={(e) => {
                e.stopPropagation();
                e.preventDefault();
                setShowDeleteConfirm(true);
              }}
            >
              <TrashIcon className='h-5 w-5' />
            </button>
          </div>{' '}
        </div>

        {/* Delete Confirmation Modal */}
        <ConfirmationModal
          title='Confirm Delete'
          message='Are you sure you want to delete this offer? This action cannot be undone.'
          confirmText='Delete'
          onConfirm={handleDelete}
          onCancel={() => setShowDeleteConfirm(false)}
          isOpen={showDeleteConfirm}
        />

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

          <div className='flex items-end pr-4'>
            <div className='text-right'>
              <p className='text-sm text-gray-600'>
                {offer.isAuction ? 'Current bid' : 'Price'}
              </p>
              <p className='text-2xl font-bold text-green-600'>
                {offer.price.toLocaleString()} zł
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
