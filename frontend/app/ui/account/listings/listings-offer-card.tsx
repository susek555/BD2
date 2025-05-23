'use client';

import { deleteListing } from '@/app/lib/api/listing/requests';
import { BaseOffer } from '@/app/lib/definitions/SaleOffer';
import ConfirmationModal from '@/app/ui/(common)/confirm-modal';
import GenericOfferCard from '@/app/ui/(offers-table)/generic-offer-card';
import { PencilIcon, TrashIcon } from '@heroicons/react/20/solid';
import { useState } from 'react';

export default function SingleListingsOffer({ offer }: { offer: BaseOffer }) {
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  const handleDelete = () => {
    deleteListing(offer.id);
    setShowDeleteConfirm(false);
  };

  const headerContent = (
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
    </div>
  );

  return (
    <>
      <GenericOfferCard offer={offer} headerContent={headerContent} />

      <ConfirmationModal
        title='Confirm Delete'
        message='Are you sure you want to delete this offer? This action cannot be undone.'
        confirmText='Delete'
        onConfirm={handleDelete}
        onCancel={() => setShowDeleteConfirm(false)}
        isOpen={showDeleteConfirm}
      />
    </>
  );
}
