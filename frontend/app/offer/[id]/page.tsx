export default async function Page(props: { params: Promise<{id: string }> }) {
    const { params } = props;
    const { id } = await params;

    return (
        <div>
            <h1>Offer ID: {id}</h1>
            {/* Add your offer details here */}
        </div>
    );
}