
export function ReviewCardSkeleton() {
  return (
    <div className="rounded-lg border border-gray-200 p-4 shadow-sm animate-pulse">
      <div className="flex justify-between items-center mb-2">
        <div className="h-4 bg-gray-200 rounded w-1/3"></div>
        <div className="h-3 bg-gray-200 rounded w-1/4"></div>
      </div>
      <div className="mb-2">
        <div className="flex items-center">
          {Array.from({ length: 5 }).map((_, i) => (
            <div key={i} className="h-4 w-4 bg-gray-200 rounded mr-1"></div>
          ))}
        </div>
      </div>
      <div className="h-4 bg-gray-200 rounded w-full mb-2"></div>
      <div className="h-4 bg-gray-200 rounded w-3/4"></div>
    </div>
  );
}

export function ReviewGridSkeleton() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      {Array.from({ length: 4 }).map((_, index) => (
        <ReviewCardSkeleton key={index} />
      ))}
    </div>
  );
}
