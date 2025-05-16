import { ReviewGrid } from '@/app/ui/review/review-grid';
import { ReviewGridSkeleton } from '@/app/ui/review/skeletons';
import { Suspense } from 'react';

export default function ReviewsPage({ userId }: { userId: number }) {
  return (
    <div className='space-y-8'>
      <div>
        <h2 className='mb-4 text-xl font-semibold'>Reviews About You</h2>
        <Suspense fallback={<ReviewGridSkeleton />}>
          <ReviewGrid variant='for' userId={userId} />
        </Suspense>
      </div>
    </div>
  );
}
