var __path__ = '../../src/components/ArtworkImage.js';

jest.dontMock(__path__);

describe('ArtworkImage', function() {
  var React, TestUtils, ArtworkImage;
  var artworkImage;

  beforeEach(function() {
    React = require('react/addons');
    TestUtils = React.addons.TestUtils;
    ArtworkImage = require(__path__);

    artworkImage = TestUtils.renderIntoDocument(
      <ArtworkImage path='/artwork/19199193' />
    )
  });

  it('works', function() {
    expect(true).toBeTruthy();
  });
});
