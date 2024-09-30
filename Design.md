# refund form
Use this example to quickly build a UI to handle refunds for your customers. Search for a user by email, render a table of their applicable charges to action, and require a confirmation before executing the refund.
```ts
import { Action, io, ctx } from "@interval/sdk";

export default new Action(async () => {
  const customerEmail = await io.input.email(
    "Email of the customer to refund:"
  );

  const charges = await getCharges(customerEmail);

  const chargesToRefund = await io.select.table(
    "Select one or more charges to refund",
    {
      data: charges,
    }
  );

  await ctx.loading.start({
    title: "Refunding charges",
    // Because we specified `itemsInQueue`, Interval will render a progress bar versus an indeterminate loading indicator.
    itemsInQueue: chargesToRefund.length,
  });

  for (const charge of chargesToRefund) {
    await refundCharge(charge.id);
    await ctx.loading.completeOne();
  }

  // Values returned from actions are automatically stored with Interval transaction logs
  return { chargesRefunded: chargesToRefund.length };
});

```

# User settings form
Interval makes it easy to build a basic form for updating any resource. This example sketches out an admin tool for managing a user's settings and information. It supports updating a variety of data types with confirmation steps and detailed context for fields, all while providing you the flexibility to customize data access controls and form logic form your use-case.

```ts
import { Action, io, ctx } from "@interval/sdk";
import { faker } from "@faker-js/faker";

const users = [...Array(10)].map((_, i) => {
  return {
    id: faker.datatype.uuid(),
    name: faker.name.findName(),
    email: faker.internet.email(),
    avatar: faker.image.avatar(),
    address: faker.address.streetAddress(),
    subscriptionPlan: faker.helpers.arrayElement(["free", "pro", "enterprise"]),
    emailVerified: faker.datatype.boolean(),
  };
});

export default new Action(async () => {
  const user = await io.search("Search for user by name or email", {
    initialResults: users,
    renderResult: user => ({
      label: user.name,
      description: user.email,
      imageUrl: user.avatar,
    }),
    onSearch: async query => {
      return users.filter(user => {
        return (
          user.name.toLowerCase().includes(query.toLowerCase()) ||
          user.email.toLowerCase().includes(query.toLowerCase())
        );
      });
    },
  });

  ctx.log(`User selected: ${user.name}`);

  const [_, name, email, subscriptionPlan, resetPassword] = await io.group([
    io.display.heading("Update user details"),
    io.input.text("Name", {
      defaultValue: user.name,
    }),
    io.input.email("Email", {
      defaultValue: user.email,
      helpText: `Existing email (${user.email}) is ${
        user.emailVerified ? "verified" : "not verified"
      }`,
    }),
    io.select.single("Subscription plan", {
      options: ["free", "pro", "enterprise"],
      defaultValue: user.subscriptionPlan,
    }),
    io.input.boolean("Reset password?", {
      helpText: "This will log the user out of all current sessions",
    }),
  ]);

  const messages = ["update the user"];
  let helpText;
  if (resetPassword) {
    messages.push("reset their password");
    helpText = "This will log the user out all current sessions";
  }
  const confirmed = await io.confirm(
    `Are you sure you want to ${messages.join(" and ")}?`,
    {
      helpText,
    }
  );

  if (!confirmed) {
    return "No confirmation, did not update the user";
  }

  const userIndex = users.findIndex(u => u.id === user.id);
  const updatedUser = {
    ...user,
    name,
    email,
    subscriptionPlan,
  };
  users[userIndex] = updatedUser;

  if (resetPassword) resetUserPassword(updatedUser);

  return updatedUser;
});
```


# QR code reader
```ts
import { Action, io } from "@interval/sdk";
import QRCode from "qrcode";

export default new Action(async () => {
const url = await io.input.url("URL for the QR code to link to", {
placeholder: "https://example.com",
});

const buffer = await QRCode.toBuffer(url.toString());

await io.display.image("Generated QR code", { buffer });
});
```


# Account migration tool

```ts
import { Action, io } from "@interval/sdk";
import { findUsers } from "../util";

export default new Action(async () => {
  const user = await io.search("Select an account", {
    onSearch: query => {
      return findUsers(query);
    },
    renderResult: u => ({
      label: `${u.firstName} ${u.lastName}`,
      description: u.email,
    }),
  });

  const videosFile = await io.input.file("Select a file", {
    allowedExtensions: [".json"],
  });

  const videos = await videosFile.json();

  await io.display.table("Videos to import", {
    data: videos,
    helpText: "Press Continue to run the import.",
  });

  // Add a confirmation step
  const confirmed = await io.confirm(`Import ${videos.length} videos?`);
  if (!confirmed) return "Action canceled, no videos imported";

  const importedVideos: Video[] = [];

  for (let i = 0; i < videos.length; i++) {
    // use our app's internal methods to create the missing inputs
    const thumbnailUrl = await generateThumbnail(videos[i].url);
    const slug = await getCollisionSafeSlug(videos[i].name);

    const createdAt = new Date(videos[i].createdAt * 1000);

    const video = await prisma.video.create({
      data: {
        title: videos[i].name,
        url: videos[i].url,
        thumbnailUrl,
        slug,
        createdAt,
        user: { connect: { id: user.id } },
      },
    });

    importedVideos.push(video);
  }

  return `Imported ${importedVideos.length} videos for ${user.email}`;
});

```



# Bulk issue editor
Bulk issue editor via the GitHub API

Interval works right in your existing codebase and integrates seamlessly with any external APIs you're using. This example demos a tool that leverages GitHub's API. It supports selecting and managing multiple issues at once within your repositories, all via a dynamic UI generated automatically by Interval.
```ts
export default new Action(async () => {
  const repos = await githubAPI(
    "GET",
    "https://api.github.com/user/repos?affiliation=owner&type=public"
  );
  const selectedRepos = await io.select.table(
    "Select a repo to edit issues for",
    {
      data: repos,
      columns: [
        "full_name",
        "created_at",
        { label: "owner", renderCell: row => row.owner.login },
      ],
      maxSelections: 1,
    }
  );

  const repo = selectedRepos[0];
  const issues = await githubAPI(
    "GET",
    `https://api.github.com/repos/${repo.full_name}/issues`
  );
  const labels = await githubAPI(
    "GET",
    `https://api.github.com/repos/${repo.full_name}/labels`
  );

  const selectedIssues = await io.select.table("Select issues to edit", {
    data: issues,
    columns: [
      "title",
      "created_at",
      { label: "user", renderCell: row => row.user.login },
      "url",
    ],
  });

  const operation = await io.select.single(
    "How do you want to edit the issues?",
    {
      options: ["Close" as const, "Set labels" as const],
    }
  );

  let body = {};
  switch (operation.value) {
    case "Close":
      body["state"] = "closed";
      break;
    case "Set labels":
      const selectedLabels = await io.select.multiple(
        "Select labels to assign",
        {
          options: labels.map(label => label.name),
        }
      );
      body["labels"] = selectedLabels.map(label => label.value);
      break;
  }

  for (const issue of selectedIssues) {
    const res = await githubAPI(
      "PATCH",
      `https://api.github.com/repos/${repo.full_name}/issues/${issue.number}`,
      body
    );
    ctx.log(res);
    ctx.loading.completeOne();
  }

  return "Done!";
});
```

# Quick and easy metrics notifier

With Interval you can:

    Send email or Slack notifications from actions
    Schedule actions to be run on a recurring basis

Since Interval already works alongside your database and backend, these features in combination make it trivial to quickly spin up reminders to keep your team aware of important metrics. You can even send notifications on a particular threshold or condition, which is useful to maintain data integrity and alert when expectations are not met.

Once the action is backgroundable, a recurring schedule can be set on the action's settings page in the Interval dashboard.

```ts
import { Action, ctx } from "@interval/sdk";
import { PrismaClient } from "@prisma/client";
const prisma = new PrismaClient();

export async function activeUsers() {
  return await prisma.$queryRaw`SELECT count(*) FROM "Users" WHERE active = true`;
}

export default new Action({
  backgroundable: true,
  handler: async () => {
    const count = await activeUsers();
    ctx.log("Active user count:", count);

    await ctx.notify({
      message: `As of today, we have *${count}* monthly active users`,
      title: "Latest active user count 📈",
      delivery: [
        {
          to: "#metrics",
          method: "SLACK",
        },
      ],
    });
  },
});
```


# Website comparison

When deploying sites and apps on the web, you'll often want to test that code updates don't break or unexpectedly change existing designs. In this example, we use Interval to build a tool that takes snapshots of any two webpages (e.g. production vs staging), ensures they look exactly the same, and highlights any differences between them.
```ts
import { Action, io, ctx } from "@interval/sdk";
import puppeteer from "puppeteer";
import { PNG } from "pngjs";
import pixelmatch from "pixelmatch";

export default new Action(async () => {
  const [_, baseUrl, compUrl] = await io.group([
    io.display.heading(
      "Select ground truth page and comparison page to compare screenshots"
    ),
    io.input.text("Ground truth URL"),
    io.input.text("Comparison URL"),
  ]);

  ctx.loading.start("Taking screenshots...");

  const browser = await puppeteer.launch();
  const page = await browser.newPage();

  await page.goto(baseUrl);
  const baseScreenshot = await page.screenshot({
    encoding: "binary",
    fullPage: true,
  });
  const baseBase64 = baseScreenshot.toString("base64");

  await page.goto(compUrl);
  const compScreenshot = await page.screenshot({
    encoding: "binary",
    fullPage: true,
  });
  const compBase64 = compScreenshot.toString("base64");

  ctx.loading.update("Comparing screenshots...");

  const basePng = PNG.sync.read(baseScreenshot);
  const compPng = PNG.sync.read(compScreenshot);
  const diff = new PNG({ width: basePng.width, height: basePng.height });

  const diffPixels = pixelmatch(
    basePng.data,
    compPng.data,
    diff.data,
    basePng.width,
    basePng.height,
    { threshold: 0.1 }
  );

  if (diffPixels > 0) {
    const diffBase64 = PNG.sync.write(diff).toString("base64");

    await io.display.markdown(
      `
## ⚠️ The pages are different

|      |      |      |
| ---- | ---- | ---- |
| ![Screenshot](data:image/png;base64,${baseBase64}) | ![Screenshot](data:image/png;base64,${diffBase64}) | ![Screenshot](data:image/png;base64,${compBase64}) |
`
    );
  } else {
    await io.display.markdown(`
## ✅ These pages are the same

![Screenshot](data:image/png;base64,${baseBase64})
`);
  }

  return "Done!";
});
```


# Concepts

## Actions
An action is a function you write that runs in your codebase and can be called from your team's Interval dashboard. Actions are the fundamental building block for creating tools with Interval.

What can you do inside an action?

    Collect input and display output with I/O methods.
    Send notifications using ctx.notify.
    Log information from using ctx.log.
    Access the information like which user is running the action using other ctx properties.



## Environments

Actions in the Interval dashboard exist within an environment. When you're using the Interval dashboard, you're always in a particular environment, and you'll only be able to view and run actions within that environment.

You can change environments using the switcher in the navigation bar.

Interval comes with two environments by default, Production and Development.

Production is for live actions to be shared with your whole organization, while Development is designed for engineers to rapidly develop actions in a local environment. Each user account with a developer or admin role is assigned its own Development environment and an API key for connecting to this environment, which is referred to as your Personal Development Key.

The environment that your Interval actions run within is determined by the API key that you provide w

## BFF

Bff is the entrypoint for defining actions and pages.

## Pages

Represent a collection of actions that are grouped together. Using pages you can organise your site by grouping actions together in pages.

## Transactions

A transaction is an instance of an action run.

For example, if your team has a Refund User action, every time someone runs that action, a transaction will be created.

You can view a history of all transactions from the History page in your team's Interval dashboard.

If you want to store information with an action along with your Transaction history, just return that information from your action.
Completion state

A completed transaction can have one of three states:

    Success occurs when a transaction returns without throwing an error.
    Error occurs when a transaction's execution is interrupted by an error being thrown.
    Canceled occurs when the person running an action ends execution early. For a backgroundable action, this happens when the person running the action explicitly cancels it. Otherwise, a transaction will appear Canceled when the person running the action leaves the page early.