-- CreateTable
CREATE TABLE "Talkshow" (
    "id" TEXT NOT NULL,
    "kanal" TEXT NOT NULL,
    "host" TEXT NOT NULL,

    CONSTRAINT "Talkshow_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "SearchQuery" (
    "id" TEXT NOT NULL,
    "active" BOOLEAN NOT NULL,
    "tags" TEXT[],
    "talkshowId" TEXT NOT NULL,

    CONSTRAINT "SearchQuery_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "SearchQuery" ADD CONSTRAINT "SearchQuery_talkshowId_fkey" FOREIGN KEY ("talkshowId") REFERENCES "Talkshow"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
