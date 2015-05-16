var rewire = require('rewire');

describe('ArtworkImage', function() {
  var React, TestUtils, ArtworkImage;
  var domNode, artworkImage;

  beforeEach(function() {
    React = require('react/addons');
    TestUtils = React.addons.TestUtils;
    ArtworkImage = rewire('../../src/components/ArtworkImage.js');

    artworkImage = TestUtils.renderIntoDocument(
      <ArtworkImage path='/artwork/19199193' />
    );
    domNode = React.findDOMNode(artworkImage);
  });

  describe('in the initial state', function() {
    it('should not be visible', function() {
      var classes = domNode.getAttribute('class').split(' ');
      expect(classes).not.to.contain('visible');
    });
  });

  describe('after the image has loaded', function() {
    beforeEach(function() {
      TestUtils.Simulate.load(domNode);
    });

    it('should be visible', function() {
      var classes = domNode.getAttribute('class').split(' ');
      expect(classes).to.contain('visible');
    });
  });

  describe('after the image has errored', function() {
    beforeEach(function() {
      TestUtils.Simulate.error(domNode);
    });

    it('should be visible', function() {
      var classes = domNode.getAttribute('class').split(' ');
      expect(classes).not.to.contain('visible');
    });
  });

});
