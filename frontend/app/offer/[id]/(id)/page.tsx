import { getMyCurrentBid } from '@/app/lib/api/offer/bid';
import { authConfig } from '@/app/lib/authConfig';
import { fetchOfferDetails } from '@/app/lib/data/offer/data';
import { fetchAverageRating } from '@/app/lib/data/reviews/data';
import OfferDescription from '@/app/ui/offer/[id]/description';
import OfferDetailsTable from '@/app/ui/offer/[id]/details-table';
import Favourite from '@/app/ui/offer/[id]/favourite';
import OwnerView from '@/app/ui/offer/[id]/owner-view';
import Photos from '@/app/ui/offer/[id]/photos';
import Price from '@/app/ui/offer/[id]/price';
import UserDetails from '@/app/ui/offer/[id]/user-details';
import { getServerSession } from 'next-auth/next';

export default async function Page(props: { params: Promise<{ id: string }> }) {
  const { params } = props;
  const { id } = await params;
  console.log(id);

  const offer = await fetchOfferDetails(id);
  if (!offer) {
    throw new Error('Offer not found');
  }

  // TODO check if price updating is working after bid or make request for highest bid

  const session = await getServerSession(authConfig);
  const isLoggedIn = !!session;
  const user_id = session?.user?.userId;
  const username = session?.user?.username;

  console.log('user_id', user_id);
  console.log("is Auction", offer.isAuction);
  console.log("condition", offer.isAuction && user_id);

  const [sellerAverageRating, myCurrentBid] = await Promise.all([
    fetchAverageRating(offer.sellerId),
    // TODO: Uncomment when getMyCurrentBid is implemented
    // (offer.isAuction && user_id) ? getMyCurrentBid(offer.id, user_id!) : Promise.resolve(undefined)
    Promise.resolve(undefined)
  ]);

  console.log('My current bid:', myCurrentBid);


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
            <div className='flex flex-row gap-5'>
              <h1 className='text-3xl font-bold'>{offer.name}</h1>
              {username !== offer.sellerName && isLoggedIn && (
                <Favourite isFavourite={offer.is_favourite} id={id} />
              )}
            </div>
          </div>
          <div className='my-10' />
          {offer.can_edit ||
          offer.can_delete ||
          username === offer.sellerName ? (
            <>
              <OwnerView
                can_edit={offer.can_edit}
                can_delete={offer.can_delete}
                offer_id={id}
                isAuction={offer.isAuction}
              />
              <div className='my-5' />
            </>
          ) : (
            <></>
          )}
          <>
            <Price
              data={{
                id: id,
                price: offer.price ?? 0,
                isAuction: offer.isAuction,
                auction: offer.auctionData,
                isActive: offer.isActive,
                myCurrentBid: myCurrentBid,
                priceOnly:
                  offer.can_edit ||
                  offer.can_delete ||
                  username === offer.sellerName,
              }}
              loggedIn={isLoggedIn}
            />
          </>

          <div className='my-4' />
          <UserDetails
            sellerId={offer.sellerId}
            sellerName={offer.sellerName}
            sellerAverageRating={sellerAverageRating}
            offerId={Number(id)}
          />
          {/* // TODO - maybe add google maps with location */}
        </div>
      </div>
    </>
  );
}
