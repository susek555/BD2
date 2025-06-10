export function ProfileInfoSkeleton() {
  return (
    <div className='rounded-lg bg-white p-6 shadow'>
      <div className='flex flex-col items-center gap-4 md:flex-row md:items-start'>
        <div className='h-24 w-24 animate-pulse rounded-full bg-gray-200' />
        <div className='flex flex-1 flex-col gap-3 text-center md:text-left'>
          <div className='h-6 w-40 animate-pulse rounded bg-gray-200' />
          <div className='h-4 w-48 animate-pulse rounded bg-gray-200' />
          <div className='h-4 w-32 animate-pulse rounded bg-gray-200' />
        </div>
      </div>
    </div>
  );
}

export function TabBarSkeleton() {
  return (
    <div className='border-b border-gray-200 bg-white'>
      <nav className='flex gap-2 px-4 md:px-6'>
        {[1, 2, 3, 4, 5].map((i) => (
          <div
            key={i}
            className='my-2 h-10 w-24 animate-pulse rounded bg-gray-200'
          />
        ))}
      </nav>
    </div>
  );
}

export function TabContentSkeleton() {
  return (
    <div className='rounded-lg bg-white shadow'>
      <div className='p-6'>
        <div className='mb-4 h-8 w-32 animate-pulse rounded bg-gray-200' />
        <div className='space-y-3'>
          <div className='h-4 w-full animate-pulse rounded bg-gray-200' />
          <div className='h-4 w-3/4 animate-pulse rounded bg-gray-200' />
          <div className='h-4 w-1/2 animate-pulse rounded bg-gray-200' />
        </div>
      </div>
    </div>
  );
}

export function ActivitySkeleton() {
  return (
    <div className='space-y-4'>
      {[1, 2, 3, 4].map((i) => (
        <div key={i} className='rounded border border-gray-200 p-4'>
          <div className='mb-2 h-4 w-1/4 animate-pulse rounded bg-gray-200' />
          <div className='h-4 w-3/4 animate-pulse rounded bg-gray-200' />
        </div>
      ))}
    </div>
  );
}

export function ListingsSkeleton() {
  return (
    <div className='grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3'>
      {[1, 2, 3, 4, 5, 6].map((i) => (
        <div key={i} className='rounded-lg border border-gray-200 p-4'>
          <div className='mb-4 h-48 animate-pulse rounded bg-gray-200' />
          <div className='mb-2 h-4 w-3/4 animate-pulse rounded bg-gray-200' />
          <div className='h-4 w-1/2 animate-pulse rounded bg-gray-200' />
        </div>
      ))}
    </div>
  );
}

export function FavoritesSkeleton() {
  return (
    <div className='grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3'>
      {[1, 2, 3, 4, 5, 6].map((i) => (
        <div key={i} className='rounded-lg border border-gray-200 p-4'>
          <div className='mb-4 h-48 animate-pulse rounded bg-gray-200' />
          <div className='mb-2 h-4 w-3/4 animate-pulse rounded bg-gray-200' />
          <div className='h-4 w-1/2 animate-pulse rounded bg-gray-200' />
        </div>
      ))}
    </div>
  );
}

export function ReviewsSkeleton() {
  return (
    <div className='space-y-4'>
      {[1, 2, 3].map((i) => (
        <div key={i} className='rounded-lg border border-gray-200 p-4'>
          <div className='mb-3 flex items-center gap-2'>
            <div className='h-5 w-20 animate-pulse rounded bg-gray-200' />
            <div className='h-4 w-32 animate-pulse rounded bg-gray-200' />
          </div>
          <div className='space-y-2'>
            <div className='h-4 w-full animate-pulse rounded bg-gray-200' />
            <div className='h-4 w-5/6 animate-pulse rounded bg-gray-200' />
          </div>
        </div>
      ))}
    </div>
  );
}

export function SettingsSkeleton() {
  return (
    <div className='space-y-6'>
      <div className='rounded-lg border border-gray-200 p-4'>
        <div className='mb-4 h-6 w-40 animate-pulse rounded bg-gray-200' />
        <div className='space-y-4'>
          <div>
            <div className='mb-2 h-4 w-24 animate-pulse rounded bg-gray-200' />
            <div className='h-10 w-full animate-pulse rounded bg-gray-200' />
          </div>
          <div>
            <div className='mb-2 h-4 w-24 animate-pulse rounded bg-gray-200' />
            <div className='h-10 w-full animate-pulse rounded bg-gray-200' />
          </div>
        </div>
      </div>
    </div>
  );
}
