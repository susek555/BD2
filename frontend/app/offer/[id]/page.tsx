import { fetchOfferDetails } from "@/app/lib/data";
import Photos from "@/app/ui/offer/[id]/photos";
import Price from "@/app/ui/offer/[id]/price";
import OfferDescription from "@/app/ui/offer/[id]/description";
import OfferDetailsTable from "@/app/ui/offer/[id]/details-table";
import UserDetails from "@/app/ui/offer/[id]/user-details";

export default async function Page(props: { params: Promise<{ id: string }> }) {
  const { params } = props;
  const { id } = await params;
  console.log(id);

  const offer = await fetchOfferDetails(id);
  if (!offer) {
    throw new Error('Offer not found');
  }

  return (
    <>
      <div className='my-3' />
      <div className='flex flex-row gap-10 md:px-20'>
        <div className='flex w-full flex-col gap-4'>
          <div className='h-full w-full border-0 border-gray-300'>
            <div className='static w-full'>
              <Photos imagesURLs={offer.imagesURLs} />
            </div>
          </div>
          <div className='my-3' />
          <OfferDescription description={offer.description} />
          <div className='my-3' />
          <OfferDetailsTable details={offer.details} />
          <div className='my-10' />
        </div>
        <div className='h-full md:h-130 md:w-120'>
          <div className='p-4'>
            <h1 className='text-3xl font-bold'>{offer.name}</h1>
          </div>
          <div className='my-10' />
          <Price
            data={{
              id: id,
              price: offer.price ?? 0,
              isAuction: offer.isAuction,
              auction: offer.auctionData,
              isActive: offer.isActive,
            }}
          />
          <div className='my-4' />
          <UserDetails sellerName={offer.sellerName} />
          {/* // TODO - maybe add google maps with location */}
        </div>
      </div>
    </>
  );
}
