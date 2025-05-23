import { UserCircleIcon } from '@heroicons/react/20/solid';
import { AverageRatingCard } from '../../review/average-rating-card';

export default function UserDetails({
  sellerName,
  sellerAverageRating,
}: {
  sellerName: string;
  sellerAverageRating: number;
}) {
  return (
    <div className='flex h-full w-full flex-col gap-4 border border-gray-300 md:h-40 md:w-120'>
      <div className='flex h-full flex-col items-center justify-center gap-2'>
        <p className='text-2xl'>Seller</p>
        <a
          href={`/account/${sellerName}`}
          // TODO - add link to user profile
          className='flex flex-row items-center justify-center rounded-md bg-gray-100 p-4 transition hover:bg-gray-200 md:w-100'
        >
          <div className='flex flex-row items-center gap-2'>
            <div className='flex h-8 w-8 items-center justify-center rounded-full bg-blue-500'>
              <UserCircleIcon className='h-8 w-8 text-white' />
            </div>
            <p className='text-2xl font-bold'>{sellerName}</p>
          </div>
          <AverageRatingCard rating={sellerAverageRating} />
        </a>
      </div>
    </div>
  );
}
