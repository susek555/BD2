'use client';

import {
  deleteListingAction,
  revalidateListingsAction,
} from '@/app/actions/listing-actions';
import { BaseOffer } from '@/app/lib/definitions/SaleOffer';
import ConfirmationModal from '@/app/ui/(common)/confirm-modal';
import GenericOfferCard from '@/app/ui/(offers-table)/generic-offer-card';
import { PencilIcon, TrashIcon } from '@heroicons/react/20/solid';
import { useState } from 'react';
import toast from 'react-hot-toast';

export default function SingleListingsOffer({ offer }: { offer: BaseOffer }) {
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  const handleDelete = async () => {
    setIsDeleting(true);

    try {
      const result = await deleteListingAction(offer.id);

      if (result.success) {
        setShowDeleteConfirm(false);
        revalidateListingsAction();
      } else {
        toast.error('Failed to delete offer. Please try again');
      }
    } catch (error) {
      toast.error('Failed to delete offer. Please try again.');
    } finally {
      setIsDeleting(false);
    }
  };
  const headerContent = (
    <div className='flex space-x-2 pr-4'>
      <button
        disabled={!offer.can_modify}
        className={`rounded-md px-3 py-1 text-sm font-medium text-white transition-colors ${
          offer.can_modify
            ? 'bg-blue-600 hover:bg-blue-700'
            : 'cursor-not-allowed bg-gray-400'
        }`}
        onClick={(e) => {
          e.stopPropagation();
          e.preventDefault();
          window.location.href = `/offer/${offer.id}/edit`;
        }}
      >
        <PencilIcon className='h-5 w-5' />
      </button>
      <button
        disabled={!offer.can_modify}
        className={`rounded-md px-3 py-1 text-sm font-medium text-white transition-colors ${
          offer.can_modify
            ? 'bg-red-600 hover:bg-red-700'
            : 'cursor-not-allowed bg-gray-400'
        }`}
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
        confirmText={isDeleting ? 'Deleting...' : 'Delete'}
        onConfirm={handleDelete}
        onCancel={() => setShowDeleteConfirm(false)}
        isOpen={showDeleteConfirm}
        disabled={isDeleting}
      />{' '}
    </>
  );
}
